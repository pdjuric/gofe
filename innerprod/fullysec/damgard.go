/*
 * Copyright (c) 2018 XLAB d.o.o
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fullysec

import (
	"fmt"
	"math/big"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/internal"
	"github.com/fentec-project/gofe/internal/dlog"
	"github.com/fentec-project/gofe/internal/keygen"
	"github.com/fentec-project/gofe/sample"
)

// DamgardParams includes public parameters for the Damgard inner
// product scheme.
// L (int): The length of vectors to be encrypted.
// Bound (int): The value by which coordinates of vectors x and y are bounded.
// G (int): Generator of a cyclic group Z_P: G**(Q) = 1 (mod P).
// H (int): Generator of a cyclic group Z_P: H**(Q) = 1 (mod P).
// P (int): Modulus - we are operating in a cyclic group Z_P.
// Q (int): Multiplicative order of G and H.
type DamgardParams struct {
	L     int
	Bound *big.Int
	G     *big.Int
	H     *big.Int
	P     *big.Int
	Q     *big.Int
}

// Damgard represents a scheme instantiated from the DDH assumption
// based on DDH variant of:
// Agrawal, Shweta, Libert, and Stehle:
// "Fully secure functional encryption for inner products,
// from standard assumptions".
type Damgard struct {
	Params *DamgardParams
}

// NewDamgard configures a new instance of the scheme.
// It accepts the length of input vectors l, the bit length of the
// modulus (we are operating in the Z_p group), and a bound by which
// coordinates of input vectors are bounded.
//
// It returns an error in case the scheme could not be properly
// configured, or if precondition l * bound² is >= order of the cyclic
// group.
func NewDamgard(l, modulusLength int, bound *big.Int) (*Damgard, error) {
	key, err := keygen.NewElGamal(modulusLength)
	if err != nil {
		return nil, err
	}
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)

	bSquared := new(big.Int).Exp(bound, two, nil)
	prod := new(big.Int).Mul(big.NewInt(int64(2*l)), bSquared)
	if prod.Cmp(key.Q) > 0 {
		return nil, fmt.Errorf("2 * l * bound^2 should be smaller than group order")
	}

	h := new(big.Int)
	for {
		sampler := sample.NewUniformRange(two, key.Q)
		r, err := sampler.Sample()
		if err != nil {
			return nil, err
		}

		// h generated in the following way is always a generator with order q
		h.Exp(key.G, r, key.P)

		// additional checks to avoid some known attacks
		if new(big.Int).Mod(new(big.Int).Sub(key.P, one), h).Cmp(zero) == 0 {
			continue
		}
		hInv := new(big.Int).ModInverse(h, key.P)
		if new(big.Int).Mod(new(big.Int).Sub(key.P, one), hInv).Cmp(zero) == 0 {
			continue
		}
		break
	}

	return &Damgard{
		Params: &DamgardParams{
			L:     l,
			Bound: bound,
			G:     key.G,
			H:     h,
			P:     key.P,
			Q:     key.Q,
		},
	}, nil
}

// NewDamgardPrecomp configures a new instance of the scheme based on
// precomputed prime numbers and generators.
// It accepts the length of input vectors l, the bit length of the
// modulus (we are operating in the Z_p group), and a bound by which
// coordinates of input vectors are bounded. The modulus length should
// be one of values 1024, 1536, 2048, 2560, 3072, or 4096. The precomputed
// prime numbers and generators were simply obtained by running NewDamgard
// function.
//
// It returns an error in case the scheme could not be properly
// configured, or if precondition l * bound² is >= order of the cyclic
// group.
func NewDamgardPrecomp(l, modulusLength int, bound *big.Int) (*Damgard, error) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	g := new(big.Int)
	h := new(big.Int)
	p := new(big.Int)

	if modulusLength == 1024 {
		g.SetString("34902160241479276675539633849372382885917193816560610471607073855548755350834003692485735908635894317735639518678334280193650806072183057417077181724192674928134805218882803812978345229222559213790765817899845072682155064387311523738581388872686127675360979304234957611566801734164757915959042140104663977828", 10)
		h.SetString("15420637599119437909472314391464117244272797175081469344935256550115202257944676665553901111021364396888996244331294748348639524453606317448549098209774792643437849812948134158610698367936877182339463164515756739059909678247548635902607622196992717889883047826579537373927574110724167611850755082659113353814", 10)
		p.SetString("166211269243229118758738154756726384542659478479960313411107431885216572625212662756677338184675400324411541201832214281445670912135683272416408753424543622705770319923251281963485084208425069817917631106045349238686234860629044433560424091289406000897029571960128048529362925472176997104870527051276406995203", 10)
	} else if modulusLength == 1536 {
		g.SetString("676416913692519694440150163403654362412279108516867264953779609011365998625435399420336578530015558254310139891236630566729665914687641028600402606957815727025192669238117788115237116562468680376464346714542467465836552396661693422160454402926392749202926871877212792118140354124110927269910674002861908621272286950597240072605316784317536178700101838123530590145680002962405974024190384775185108002307650499125333676880320808656556635493186351335151559453463208", 10)
		h.SetString("162435356591098080112923949933138473281129080012462876365242296621510003247926800168810467125425312741356845555051398168995345221821254788094709212569083281560861934370026731898947082535853089500762419810631794616608702382341836751059467526394750733484031555507059435514454850206195181859188711511695415306637351122791298103837179337333322817996915750092675310561842474481585085713281105392293555407844591189296291324552036671109998940504540100161197823599243891", 10)
		p.SetString("1851297899986638926486011430658634631676522135433726749065856232802142091866650774719879427474637700607873256035038534449089405369134066444876856913629831069906506096279113968447116822488133963417347136141052507685108634240736100862550194947326287783557220764070479431781692630708747550712729778398000353165406458520850089303530985563143326919073190605085889925484113854496074216626577246143598303709289292397203458923541841135799203967503522114881404128535647507", 10)
	} else if modulusLength == 2048 {
		g.SetString("4006960929413042209594165215465319088439374252008797022450541422457034721098828647078778605657155669917104962611933792130890703423519992986737966991597160684973795472419962788730248050852176194215699504914899438223683843401963466624139534923052671383315398134823370041633710463630745156269175253639670460050105594663691338308037509280576148624454011047879615100156717631945194107791315234171086603775159708325087759679758438868772220133433497821899045165244202228696902434100209752952701657306825368599999359102329396520012735146260911352901326915877502873633420811221206110021993351144711002138373506576799781061829", 10)
		h.SetString("16531397285674350767314702525644915041456240916042269006661016218466751820121367708744665988146094312460459686332899317211264306268688887145182608791297975910105233621322826639428591684097281721326762853220049531036892964645338402521408932284623118973840687103995270506055864839392780688919921574729334863152272712825111303515621796552988450890438067795005159750945818303100847992123152427222482913524403644261704298722743429153705746979942035168013929376597721413830914117313871470498590418693493905967182265847426629997026054543887503362398742858491325979767814095393383666927420156454285653125386823358683247619671", 10)
		p.SetString("28884237504713658990682089080899862128005980675308910325841161962760155725037929764087367167449843609136681034352509183117742758446654629096509285354423361556493020266963222508540306384896802796001914743293196010488452478370041404523014215612960481024232879327123268440037633547483165934132901270561772860319969916305482525766132307669097012989986613879246932730824899649301621408341438037745468033187743673001187803377254713546325789438300798311106106322698517805307792059495696632070953526611920926003483451787562399452650878943515646786958216714025307572678422373120397225912926110031401983688860264234966561627699", 10)
	} else if modulusLength == 2560 {
		g.SetString("283408881721750179985507845260248881237898607313021593637913985490973593382472548378053368228310040484438920416918021737085067966444840449068073874437662089659563355479608930182263107110969154912883370207053052795964658868443319273167870311282045348396320159912742394374241864712808383029954025256232806465551969466207671603658677963161975454703127476120201164519187150268352527923664649275471494757270139533433456630363925187498055211365480086561354743681517539297815712218419607006668655891574362066382949706266666189227897710299445185100212256741698216505337617571970963008519334554537811591236478130526432239803909461119767954934793813410765013072006162612226471775059215628326278458577643374735250370115470812597459244082296191871275203831471332697557979904062571849", 10)
		h.SetString("60603226626244984946015074611399313803245697048858823305998791050960516381707120286785720252208953282756722645401890645185550185014757676697220190895358672768431707069471110996253922983573984124316810681408685350886144002884260922773441021522086311309748129811838078781339054390074111039635714653119513182362870864625376495724028601281359016299631376501745572486338940082085075668612842777510700862160795716160806125722301037965086971367697165162183431144060195434995026056781565232838274115314469024436666926374971029865332455354830888647688402585402990945621351542202067894748432776518006655843916388537951718725340446233164828852013369347111758769617449735121060447630493643403400643959220152005725907208460642317320745167072023561550218217956652950353398499254599742", 10)
		p.SetString("403126381462544353337185732949672984701995200694926494258456907009253008321275627278199160008680481542251933923210212978304159426174307419863781623411302777276318286800690060166638633211627521530173324015527027650548199185539205697958056639406116068885574865579676651743820636201007864067569576455725489531113260031526605827601510665037511961715114944815619491261828558745083975042855063688267346905844510423020844412350570902289599734320004108557140241966071165594059732527795488131297017205383953304055105007982366596746708951250486384299368612656872813778220074826250625689603663742175288397398948456522281031888042417278385238985218731264092285591879578299600853004336936458454638992426900228708418575870946630137618851131144232868141478901063119847104013555395370887", 10)
	} else if modulusLength == 3072 {
		g.SetString("3696261717511684685041982088526264061294500298114921057816954458306445697150348237912729036967670872345042594238223317055749478029025374644864924550052402546275985983344583674703146236623453822520422465163020824494790581472736649085281450730460445260696087738043872307629635997875332076478424042345012004769107421873566499123042621978973433575500345010912635742477932006291250637245855027695943163956584173316781442078828050076620331751405548730676363847526959436516279320074682721438642683731766682502490935962962293815202487144775533102010333956118641968798500514719248831145108532912211817219793191951880318961073149276914867129023978524587935704313755469570162971499124682746476415187933097132047611840762510892175328320025164466873845777990557296853549970943298347924080102740724512079409979152285019931666423541870247789529268168448010024121369388707140296446100906359619586133848407970098685310317291828335700424602208", 10)
		h.SetString("567649039783047480227338473928434756096941793851205367034019596755160950347349183458824085567577717539018304618223485608741105961797110596645149404693961630602136531678036354529574048049387247758318737991227930011620451360293879761672154747784975220369021073972313352929486325631485787131058176462649454317449501920846554754923439009414476391939400343349482311343015897922534760972817588470659148393369453541539857229554365836585067664597417625404552690944486027240221141785988408230137708791844568226402922187996696503785944931820332617839757682381094642917330498053626106798012177968190876620737530689308284333593333125689734905825145119037797011061634852469676220493758348073580817216823443614368951436342741042435481507815482798733572408751161565276006768858533125203508129436766572987021821650216162035709837972268600512594452400036134765839973589301700524293661567605165762293059122431279400652446300525682001019196136", 10)
		p.SetString("4387756306134544957818467663802660683665166110605728231080818705443663402154316615145921798856363268744945754470238000282108344905251127487705736550297997444150840902348669718478564904142834154197029830975532074167513046443903186309497214496864577129616824062991068960005865144004932069025136224356325248036029606434443391988386519658751798077031844645051726026696307027395796695909035405241040411794836124123435225690961994089776517262574417789067836840997650095451062948856617211542724543995145259735683916440579956961657374517806591607068842498749297993409884001044324428640569001916341503645559748760311343179943896427393009949062735145363544745972252566600994034655540841225414736222780096833045470605544717177880459300618917961703559234544541206877026518430276932498602360341258899345739335298856394124351357206871568254540730107127298623178526868418799471896060015463201459762913197633841160710893895836663035998106119", 10)
	} else if modulusLength == 4096 {
		g.SetString("51665588681810560577916524923861643358980285220048008212528567741884121491554604183472728540139463099618903178110360757930742372390027135064809646425064896539133721148335557788263239281487173350543811713890328584918216783142094297306639941000480756707312457878765754357205186485080839623690156744636468433787780205323460166423447602447200754978133176713947189000663528355089645281397174452923418212485422962705227706103188302892660448134233848971142570881089940852441776074246332915421265800026335300100610273942459340241610730244726628211914068945587128124478812632725838440727321816905181830592204023095726270782834020990986443265625389712733369116937470448592846480352222814297792606318850361699893703272484112273500581408730519942517586496563772194165844831300501908379990979449691597045730512107756238377635183257797115883839801779086058652272455400286891699445584526719648220045380141260347316315487340493029966105973850214850475440630205768783542021741101804842248602349004364816943429122368563644935802417389995380389429997320053299323220481603252879925927515844929958940305561718295197935926645561977544440676439150126025681320050786964708227836328341875446457912905977470123640014345655062829575775837287500880054558386787", 10)
		h.SetString("653735051824112393106801492829043654983094810153936916566419297234893790474265474662670688710073348877202869152643983800310336936407217572878666006904653566429370273243928900524622200126090143285935697025102272467307580256707629199315519492688020243538476473287317884073943557740453015512750672402070439254723909940895097629630953500504933411614006835509266155433485165683944804036549172274160941459059251790260424550107135346982158930124883878553492540634625536558325828876220292338217093182644187426663586499630969144324838282516931412124575753931891717246058395623268249929461243359978135742469252500443989584744702727780893623276815417727584978028391169715443458364712495037724857002997953369015075369014382917440895283545982020768248209480534141352756164681807291197719347396848885385564846533227145898299352065862788210851139899108807494649254246247909836864209696565581846420881131639054554828603682225473062422081030815891806059849482326467150974089876248055494734946134244334370720707623963134127685862797676223995758737312498339183154611067568296397353652230198967761004520297654249905695573112921297631619886554083093006856350790194488256889122557656064371766221432744243315870629759094828247037312561620044735468234281300", 10)
		p.SetString("1022249395832567838406986294560330159176972202126664245047364146720891252715766488477689126342364655087193411078517616569887825896401401223927363505007778278205623713273194552498760148834874746839752870298152746450585455651115247220867383465863156721401567161663838310658875672995951663020449772454232797368263754624173026584111779206080723120076751471597509403139249260220696195263597156452889920392585797464801375940661326779247976331028637271512085826066667631423502199894046717721786935806581428328491087482664043743281068318459302242239861275878019857365021173868449409246193470959347916848019032536247915451026158871684654213802886886213841729258073333569276986893577214659899227179735448593265633219968622571880602115519942763955551007919826002851866939641065270816032435114864853636918330698605282572789904941484540512478406984407320963402583009124880812235841866246441862987563989772424040933513333746472128494254253767426962063553015635240386636751473945937412527996558505231385625318878887383161350102080329744822052478052004574860361461762694379860797225344866320388590336321515376486033237159694567932935601775209663052272120524337888258857351777348841323194553467226791591208931619058871750498804369190487499494069660723", 10)
	} else {
		return nil, fmt.Errorf("modulus length should be one of values 1024, 1536, 2048, 2560, 3072, or 4096")
	}

	q := new(big.Int).Sub(p, one)
	q.Div(q, two)

	bSquared := new(big.Int).Exp(bound, two, nil)
	prod := new(big.Int).Mul(big.NewInt(int64(2*l)), bSquared)
	if prod.Cmp(q) > 0 {
		return nil, fmt.Errorf("2 * l * bound^2 should be smaller than group order")
	}

	return &Damgard{
		Params: &DamgardParams{
			L:     l,
			Bound: bound,
			G:     g,
			H:     h,
			P:     p,
			Q:     q,
		},
	}, nil
}

// NewDamgardFromParams takes configuration parameters of an existing
// Damgard scheme instance, and reconstructs the scheme with same configuration
// parameters. It returns a new Damgard instance.
func NewDamgardFromParams(params *DamgardParams) *Damgard {
	return &Damgard{
		Params: params,
	}
}

// DamgardSecKey is a secret key for Damgard scheme.
type DamgardSecKey struct {
	S data.Vector
	T data.Vector
}

// GenerateMasterKeys generates a master secret key and master
// public key for the scheme. It returns an error in case master keys
// could not be generated.
func (d *Damgard) GenerateMasterKeys() (*DamgardSecKey, data.Vector, error) {
	// both part of masterSecretKey
	mskS := make(data.Vector, d.Params.L)
	mskT := make(data.Vector, d.Params.L)

	masterPubKey := make([]*big.Int, d.Params.L)
	sampler := sample.NewUniformRange(big.NewInt(2), d.Params.Q)

	for i := 0; i < d.Params.L; i++ {
		s, err := sampler.Sample()
		if err != nil {
			return nil, nil, err
		}
		mskS[i] = s

		t, err := sampler.Sample()
		if err != nil {
			return nil, nil, err
		}
		mskT[i] = t

		y1 := new(big.Int).Exp(d.Params.G, s, d.Params.P)
		y2 := new(big.Int).Exp(d.Params.H, t, d.Params.P)

		masterPubKey[i] = new(big.Int).Mod(new(big.Int).Mul(y1, y2), d.Params.P)

	}

	return &DamgardSecKey{S: mskS, T: mskT}, masterPubKey, nil
}

// DamgardDerivedKey is a functional encryption key for Damgard scheme.
type DamgardDerivedKey struct {
	Key1 *big.Int
	Key2 *big.Int
}

// DeriveKey takes master secret key and input vector y, and returns the
// functional encryption key. In case the key could not be derived, it
// returns an error.
func (d *Damgard) DeriveKey(masterSecKey *DamgardSecKey, y data.Vector) (*DamgardDerivedKey, error) {
	if err := y.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	key1, err := masterSecKey.S.Dot(y)
	if err != nil {
		return nil, err
	}

	key2, err := masterSecKey.T.Dot(y)
	if err != nil {
		return nil, err
	}

	k1 := new(big.Int).Mod(key1, d.Params.Q)
	k2 := new(big.Int).Mod(key2, d.Params.Q)

	return &DamgardDerivedKey{Key1: k1, Key2: k2}, nil
}

// Encrypt encrypts input vector x with the provided master public key.
// It returns a ciphertext vector. If encryption failed, error is returned.
func (d *Damgard) Encrypt(x, masterPubKey data.Vector) (data.Vector, error) {
	if err := x.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	sampler := sample.NewUniformRange(big.NewInt(2), d.Params.Q)
	r, err := sampler.Sample()
	if err != nil {
		return nil, err
	}

	ciphertext := make([]*big.Int, len(x)+2)
	// c = g^r
	// dd = h^r
	c := new(big.Int).Exp(d.Params.G, r, d.Params.P)
	ciphertext[0] = c
	dd := new(big.Int).Exp(d.Params.H, r, d.Params.P)
	ciphertext[1] = dd

	for i := 0; i < len(x); i++ {
		// e_i = h_i^r * g^x_i
		// e_i = mpk[i]^r * g^x_i
		t1 := new(big.Int).Exp(masterPubKey[i], r, d.Params.P)
		t2 := internal.ModExp(d.Params.G, x[i], d.Params.P)
		ct := new(big.Int).Mod(new(big.Int).Mul(t1, t2), d.Params.P)
		ciphertext[i+2] = ct
	}

	return data.NewVector(ciphertext), nil
}

// Decrypt accepts the encrypted vector, functional encryption key, and
// a plaintext vector y. It returns the inner product of x and y.
// If decryption failed, error is returned.
func (d *Damgard) Decrypt(cipher data.Vector, key *DamgardDerivedKey, y data.Vector) (*big.Int, error) {
	if err := y.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	num := big.NewInt(1)
	for i, ct := range cipher[2:] {
		t1 := internal.ModExp(ct, y[i], d.Params.P)
		num = num.Mod(new(big.Int).Mul(num, t1), d.Params.P)
	}

	t1 := new(big.Int).Exp(cipher[0], key.Key1, d.Params.P)
	t2 := new(big.Int).Exp(cipher[1], key.Key2, d.Params.P)

	denom := new(big.Int).Mod(new(big.Int).Mul(t1, t2), d.Params.P)
	denomInv := new(big.Int).ModInverse(denom, d.Params.P)
	r := new(big.Int).Mod(new(big.Int).Mul(num, denomInv), d.Params.P)

	bSquared := new(big.Int).Exp(d.Params.Bound, big.NewInt(2), big.NewInt(0))
	bound := new(big.Int).Mul(big.NewInt(int64(d.Params.L)), bSquared)

	calc, err := dlog.NewCalc().InZp(d.Params.P, d.Params.Q)
	if err != nil {
		return nil, err
	}
	calc = calc.WithNeg()

	res, err := calc.WithBound(bound).BabyStepGiantStep(r, d.Params.G)
	return res, err
}
