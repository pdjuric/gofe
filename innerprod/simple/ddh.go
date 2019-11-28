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

package simple

import (
	"fmt"
	"math/big"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/internal"
	"github.com/fentec-project/gofe/internal/dlog"
	"github.com/fentec-project/gofe/internal/keygen"
	"github.com/fentec-project/gofe/sample"
)

// DDHParams represents configuration parameters for the DDH scheme instance.
type DDHParams struct {
	// length of input vectors x and y
	L int
	// The value by which coordinates of input vectors x and y are bounded.
	Bound *big.Int
	// Generator of a cyclic group Z_P: G^(Q) = 1 (mod P).
	G *big.Int
	// Modulus - we are operating in a cyclic group Z_P.
	P *big.Int
	// Order of the generator G.
	Q *big.Int
}

// DDH represents a scheme instantiated from the DDH assumption,
// based on the DDH variant by
// Abdalla, Bourse, De Caro, and Pointchev:
// "Simple Functional Encryption Schemes for Inner Products".
type DDH struct {
	Params *DDHParams
}

// NewDDH configures a new instance of the scheme.
// It accepts the length of input vectors l, the bit length of the
// modulus (we are operating in the Z_p group), and a bound by which
// coordinates of input vectors are bounded.
//
// It returns an error in case the scheme could not be properly
// configured, or if precondition l * bound² is >= order of the cyclic
// group.
func NewDDH(l, modulusLength int, bound *big.Int) (*DDH, error) {
	key, err := keygen.NewElGamal(modulusLength)
	if err != nil {
		return nil, err
	}

	if new(big.Int).Mul(big.NewInt(int64(2*l)), new(big.Int).Exp(bound, big.NewInt(2), big.NewInt(0))).Cmp(key.Q) > 0 {
		return nil, fmt.Errorf("2 * l * bound^2 should be smaller than group order")
	}

	sip := DDH{
		Params: &DDHParams{
			L:     l,
			Bound: bound,
			G:     key.G,
			P:     key.P,
			Q:     key.Q,
		},
	}

	return &sip, nil
}

// NewDDHPrecomp configures a new instance of the scheme based on
// precomputed prime numbers and generators.
// It accepts the length of input vectors l, the bit length of the
// modulus (we are operating in the Z_p group), and a bound by which
// coordinates of input vectors are bounded. The modulus length should
// be one of values 1024, 1536, 2048, 2560, 3072, or 4096.
//
// It returns an error in case the scheme could not be properly
// configured, or if precondition l * bound² is >= order of the cyclic
// group.
func NewDDHPrecomp(l, modulusLength int, bound *big.Int) (*DDH, error) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)

	g := new(big.Int)
	p := new(big.Int)

	if modulusLength == 1024 {
		g.SetString("34902160241479276675539633849372382885917193816560610471607073855548755350834003692485735908635894317735639518678334280193650806072183057417077181724192674928134805218882803812978345229222559213790765817899845072682155064387311523738581388872686127675360979304234957611566801734164757915959042140104663977828", 10)
		p.SetString("166211269243229118758738154756726384542659478479960313411107431885216572625212662756677338184675400324411541201832214281445670912135683272416408753424543622705770319923251281963485084208425069817917631106045349238686234860629044433560424091289406000897029571960128048529362925472176997104870527051276406995203", 10)
	} else if modulusLength == 1536 {
		g.SetString("676416913692519694440150163403654362412279108516867264953779609011365998625435399420336578530015558254310139891236630566729665914687641028600402606957815727025192669238117788115237116562468680376464346714542467465836552396661693422160454402926392749202926871877212792118140354124110927269910674002861908621272286950597240072605316784317536178700101838123530590145680002962405974024190384775185108002307650499125333676880320808656556635493186351335151559453463208", 10)
		p.SetString("1851297899986638926486011430658634631676522135433726749065856232802142091866650774719879427474637700607873256035038534449089405369134066444876856913629831069906506096279113968447116822488133963417347136141052507685108634240736100862550194947326287783557220764070479431781692630708747550712729778398000353165406458520850089303530985563143326919073190605085889925484113854496074216626577246143598303709289292397203458923541841135799203967503522114881404128535647507", 10)
	} else if modulusLength == 2048 {
		g.SetString("4006960929413042209594165215465319088439374252008797022450541422457034721098828647078778605657155669917104962611933792130890703423519992986737966991597160684973795472419962788730248050852176194215699504914899438223683843401963466624139534923052671383315398134823370041633710463630745156269175253639670460050105594663691338308037509280576148624454011047879615100156717631945194107791315234171086603775159708325087759679758438868772220133433497821899045165244202228696902434100209752952701657306825368599999359102329396520012735146260911352901326915877502873633420811221206110021993351144711002138373506576799781061829", 10)
		p.SetString("28884237504713658990682089080899862128005980675308910325841161962760155725037929764087367167449843609136681034352509183117742758446654629096509285354423361556493020266963222508540306384896802796001914743293196010488452478370041404523014215612960481024232879327123268440037633547483165934132901270561772860319969916305482525766132307669097012989986613879246932730824899649301621408341438037745468033187743673001187803377254713546325789438300798311106106322698517805307792059495696632070953526611920926003483451787562399452650878943515646786958216714025307572678422373120397225912926110031401983688860264234966561627699", 10)
	} else if modulusLength == 2560 {
		g.SetString("283408881721750179985507845260248881237898607313021593637913985490973593382472548378053368228310040484438920416918021737085067966444840449068073874437662089659563355479608930182263107110969154912883370207053052795964658868443319273167870311282045348396320159912742394374241864712808383029954025256232806465551969466207671603658677963161975454703127476120201164519187150268352527923664649275471494757270139533433456630363925187498055211365480086561354743681517539297815712218419607006668655891574362066382949706266666189227897710299445185100212256741698216505337617571970963008519334554537811591236478130526432239803909461119767954934793813410765013072006162612226471775059215628326278458577643374735250370115470812597459244082296191871275203831471332697557979904062571849", 10)
		p.SetString("403126381462544353337185732949672984701995200694926494258456907009253008321275627278199160008680481542251933923210212978304159426174307419863781623411302777276318286800690060166638633211627521530173324015527027650548199185539205697958056639406116068885574865579676651743820636201007864067569576455725489531113260031526605827601510665037511961715114944815619491261828558745083975042855063688267346905844510423020844412350570902289599734320004108557140241966071165594059732527795488131297017205383953304055105007982366596746708951250486384299368612656872813778220074826250625689603663742175288397398948456522281031888042417278385238985218731264092285591879578299600853004336936458454638992426900228708418575870946630137618851131144232868141478901063119847104013555395370887", 10)
	} else if modulusLength == 3072 {
		g.SetString("3696261717511684685041982088526264061294500298114921057816954458306445697150348237912729036967670872345042594238223317055749478029025374644864924550052402546275985983344583674703146236623453822520422465163020824494790581472736649085281450730460445260696087738043872307629635997875332076478424042345012004769107421873566499123042621978973433575500345010912635742477932006291250637245855027695943163956584173316781442078828050076620331751405548730676363847526959436516279320074682721438642683731766682502490935962962293815202487144775533102010333956118641968798500514719248831145108532912211817219793191951880318961073149276914867129023978524587935704313755469570162971499124682746476415187933097132047611840762510892175328320025164466873845777990557296853549970943298347924080102740724512079409979152285019931666423541870247789529268168448010024121369388707140296446100906359619586133848407970098685310317291828335700424602208", 10)
		p.SetString("4387756306134544957818467663802660683665166110605728231080818705443663402154316615145921798856363268744945754470238000282108344905251127487705736550297997444150840902348669718478564904142834154197029830975532074167513046443903186309497214496864577129616824062991068960005865144004932069025136224356325248036029606434443391988386519658751798077031844645051726026696307027395796695909035405241040411794836124123435225690961994089776517262574417789067836840997650095451062948856617211542724543995145259735683916440579956961657374517806591607068842498749297993409884001044324428640569001916341503645559748760311343179943896427393009949062735145363544745972252566600994034655540841225414736222780096833045470605544717177880459300618917961703559234544541206877026518430276932498602360341258899345739335298856394124351357206871568254540730107127298623178526868418799471896060015463201459762913197633841160710893895836663035998106119", 10)
	} else if modulusLength == 4096 {
		g.SetString("51665588681810560577916524923861643358980285220048008212528567741884121491554604183472728540139463099618903178110360757930742372390027135064809646425064896539133721148335557788263239281487173350543811713890328584918216783142094297306639941000480756707312457878765754357205186485080839623690156744636468433787780205323460166423447602447200754978133176713947189000663528355089645281397174452923418212485422962705227706103188302892660448134233848971142570881089940852441776074246332915421265800026335300100610273942459340241610730244726628211914068945587128124478812632725838440727321816905181830592204023095726270782834020990986443265625389712733369116937470448592846480352222814297792606318850361699893703272484112273500581408730519942517586496563772194165844831300501908379990979449691597045730512107756238377635183257797115883839801779086058652272455400286891699445584526719648220045380141260347316315487340493029966105973850214850475440630205768783542021741101804842248602349004364816943429122368563644935802417389995380389429997320053299323220481603252879925927515844929958940305561718295197935926645561977544440676439150126025681320050786964708227836328341875446457912905977470123640014345655062829575775837287500880054558386787", 10)
		p.SetString("1022249395832567838406986294560330159176972202126664245047364146720891252715766488477689126342364655087193411078517616569887825896401401223927363505007778278205623713273194552498760148834874746839752870298152746450585455651115247220867383465863156721401567161663838310658875672995951663020449772454232797368263754624173026584111779206080723120076751471597509403139249260220696195263597156452889920392585797464801375940661326779247976331028637271512085826066667631423502199894046717721786935806581428328491087482664043743281068318459302242239861275878019857365021173868449409246193470959347916848019032536247915451026158871684654213802886886213841729258073333569276986893577214659899227179735448593265633219968622571880602115519942763955551007919826002851866939641065270816032435114864853636918330698605282572789904941484540512478406984407320963402583009124880812235841866246441862987563989772424040933513333746472128494254253767426962063553015635240386636751473945937412527996558505231385625318878887383161350102080329744822052478052004574860361461762694379860797225344866320388590336321515376486033237159694567932935601775209663052272120524337888258857351777348841323194553467226791591208931619058871750498804369190487499494069660723", 10)
	} else {
		return nil, fmt.Errorf("modulus length should be one of values 1024, 1536, 2048, 2560, 3072, or 4096")
	}

	q := new(big.Int).Sub(p, one)
	q.Div(q, two)

	if new(big.Int).Mul(big.NewInt(int64(2*l)), new(big.Int).Exp(bound, two, zero)).Cmp(q) > 0 {
		return nil, fmt.Errorf("2 * l * bound^2 should be smaller than group order")
	}

	sip := DDH{
		Params: &DDHParams{
			L:     l,
			Bound: bound,
			G:     g,
			P:     p,
			Q:     q,
		},
	}

	return &sip, nil
}

// NewDDHFromParams takes configuration parameters of an existing
// DDH scheme instance, and reconstructs the scheme with same configuration
// parameters. It returns a new DDH instance.
func NewDDHFromParams(params *DDHParams) *DDH {
	return &DDH{
		Params: params,
	}
}

// GenerateMasterKeys generates a pair of master secret key and master
// public key for the scheme. It returns an error in case master keys
// could not be generated.
func (d *DDH) GenerateMasterKeys() (data.Vector, data.Vector, error) {
	masterSecKey := make(data.Vector, d.Params.L)
	masterPubKey := make(data.Vector, d.Params.L)
	sampler := sample.NewUniformRange(big.NewInt(2), d.Params.Q)

	for i := 0; i < d.Params.L; i++ {
		x, err := sampler.Sample()
		if err != nil {
			return nil, nil, err
		}
		y := internal.ModExp(d.Params.G, x, d.Params.P)
		masterSecKey[i] = x
		masterPubKey[i] = y
	}

	return masterSecKey, masterPubKey, nil
}

// DeriveKey takes master secret key and input vector y, and returns the
// functional encryption key. In case the key could not be derived, it
// returns an error.
func (d *DDH) DeriveKey(masterSecKey, y data.Vector) (*big.Int, error) {
	if err := y.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	key, err := masterSecKey.Dot(y)
	if err != nil {
		return nil, err
	}
	return new(big.Int).Mod(key, d.Params.Q), nil
}

// Encrypt encrypts input vector x with the provided master public key.
// It returns a ciphertext vector. If encryption failed, error is returned.
func (d *DDH) Encrypt(x, masterPubKey data.Vector) (data.Vector, error) {
	if err := x.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	sampler := sample.NewUniformRange(big.NewInt(2), d.Params.Q)
	r, err := sampler.Sample()
	if err != nil {
		return nil, err
	}

	ciphertext := make([]*big.Int, len(x)+1)
	// ct0 = g^r
	ct0 := new(big.Int).Exp(d.Params.G, r, d.Params.P)
	ciphertext[0] = ct0

	for i := 0; i < len(x); i++ {
		// ct_i = h_i^r * g^x_i
		// ct_i = mpk[i]^r * g^x_i
		t1 := new(big.Int).Exp(masterPubKey[i], r, d.Params.P)
		t2 := internal.ModExp(d.Params.G, x[i], d.Params.P)
		ct := new(big.Int).Mod(new(big.Int).Mul(t1, t2), d.Params.P)
		ciphertext[i+1] = ct
	}

	return ciphertext, nil
}

// Decrypt accepts the encrypted vector, functional encryption key, and
// a plaintext vector y. It returns the inner product of x and y.
// If decryption failed, error is returned.
func (d *DDH) Decrypt(cipher data.Vector, key *big.Int, y data.Vector) (*big.Int, error) {
	if err := y.CheckBound(d.Params.Bound); err != nil {
		return nil, err
	}

	num := big.NewInt(1)
	for i, ct := range cipher[1:] {
		t1 := internal.ModExp(ct, y[i], d.Params.P)
		num = num.Mod(new(big.Int).Mul(num, t1), d.Params.P)
	}

	denom := internal.ModExp(cipher[0], key, d.Params.P)
	denomInv := new(big.Int).ModInverse(denom, d.Params.P)
	r := new(big.Int).Mod(new(big.Int).Mul(num, denomInv), d.Params.P)

	bound := new(big.Int).Mul(big.NewInt(int64(d.Params.L)), new(big.Int).Exp(d.Params.Bound, big.NewInt(2), big.NewInt(0)))

	calc, err := dlog.NewCalc().InZp(d.Params.P, d.Params.Q)
	if err != nil {
		return nil, err
	}
	calc = calc.WithNeg()

	res, err := calc.WithBound(bound).BabyStepGiantStep(r, d.Params.G)
	return res, err

}
