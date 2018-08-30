package keygen

import (
	"fmt"
	logger "log"
	"os"
	"strings"

	"github.com/snail007/goproxy/services"
	"github.com/snail007/goproxy/utils"
	"github.com/snail007/goproxy/utils/cert"
)

type KeygenArgs struct {
	CaName     *string
	CertName   *string
	Sign       *bool
	SignDays   *int
	CommonName *string
}

type Keygen struct {
	cfg KeygenArgs
	log *logger.Logger
}

func NewKeygen() services.Service {
	return &Keygen{}
}
func (s *Keygen) CheckArgs() (err error) {
	if *s.cfg.Sign && (*s.cfg.CertName == "" || *s.cfg.CaName == "") {
		err = fmt.Errorf("ca name and cert name required for signin")
		return
	}
	if !*s.cfg.Sign && *s.cfg.CaName == "" {
		err = fmt.Errorf("ca name required")
		return
	}
	if *s.cfg.CommonName == "" {
		domainSubfixList := []string{".com", ".edu", ".gov", ".int", ".mil", ".net", ".org", ".biz", ".info", ".pro", ".name", ".museum", ".coop", ".aero", ".xxx", ".idv", ".ac", ".ad", ".ae", ".af", ".ag", ".ai", ".al", ".am", ".an", ".ao", ".aq", ".ar", ".as", ".at", ".au", ".aw", ".az", ".ba", ".bb", ".bd", ".be", ".bf", ".bg", ".bh", ".bi", ".bj", ".bm", ".bn", ".bo", ".br", ".bs", ".bt", ".bv", ".bw", ".by", ".bz", ".ca", ".cc", ".cd", ".cf", ".cg", ".ch", ".ci", ".ck", ".cl", ".cm", ".cn", ".co", ".cr", ".cu", ".cv", ".cx", ".cy", ".cz", ".de", ".dj", ".dk", ".dm", ".do", ".dz", ".ec", ".ee", ".eg", ".eh", ".er", ".es", ".et", ".eu", ".fi", ".fj", ".fk", ".fm", ".fo", ".fr", ".ga", ".gd", ".ge", ".gf", ".gg", ".gh", ".gi", ".gl", ".gm", ".gn", ".gp", ".gq", ".gr", ".gs", ".gt", ".gu", ".gw", ".gy", ".hk", ".hm", ".hn", ".hr", ".ht", ".hu", ".id", ".ie", ".il", ".im", ".in", ".io", ".iq", ".ir", ".is", ".it", ".je", ".jm", ".jo", ".jp", ".ke", ".kg", ".kh", ".ki", ".km", ".kn", ".kp", ".kr", ".kw", ".ky", ".kz", ".la", ".lb", ".lc", ".li", ".lk", ".lr", ".ls", ".lt", ".lu", ".lv", ".ly", ".ma", ".mc", ".md", ".mg", ".mh", ".mk", ".ml", ".mm", ".mn", ".mo", ".mp", ".mq", ".mr", ".ms", ".mt", ".mu", ".mv", ".mw", ".mx", ".my", ".mz", ".na", ".nc", ".ne", ".nf", ".ng", ".ni", ".nl", ".no", ".np", ".nr", ".nu", ".nz", ".om", ".pa", ".pe", ".pf", ".pg", ".ph", ".pk", ".pl", ".pm", ".pn", ".pr", ".ps", ".pt", ".pw", ".py", ".qa", ".re", ".ro", ".ru", ".rw", ".sa", ".sb", ".sc", ".sd", ".se", ".sg", ".sh", ".si", ".sj", ".sk", ".sl", ".sm", ".sn", ".so", ".sr", ".st", ".sv", ".sy", ".sz", ".tc", ".td", ".tf", ".tg", ".th", ".tj", ".tk", ".tl", ".tm", ".tn", ".to", ".tp", ".tr", ".tt", ".tv", ".tw", ".tz", ".ua", ".ug", ".uk", ".um", ".us", ".uy", ".uz", ".va", ".vc", ".ve", ".vg", ".vi", ".vn", ".vu", ".wf", ".ws", ".ye", ".yt", ".yu", ".yr", ".za", ".zm", ".zw"}
		CN := strings.ToLower(utils.RandString(int(utils.RandInt(4)%10)) + domainSubfixList[int(utils.RandInt(4))%len(domainSubfixList)])
		*s.cfg.CommonName = CN
	}
	return
}
func (s *Keygen) Start(args interface{}, log *logger.Logger) (err error) {
	s.log = log
	s.cfg = args.(KeygenArgs)
	if err = s.CheckArgs(); err != nil {
		return
	}
	if *s.cfg.Sign {
		caCert, caKey, err := cert.ParseCertAndKey(*s.cfg.CaName+".crt", *s.cfg.CaName+".key")
		if err != nil {
			return err
		}
		err = cert.CreateSignCertToFile(caCert, caKey, *s.cfg.CommonName, *s.cfg.SignDays, *s.cfg.CertName)
	} else {
		err = cert.CreateCaToFile(*s.cfg.CaName, *s.cfg.CommonName, *s.cfg.SignDays)
	}
	if err != nil {
		return
	}
	s.log.Println("success")
	os.Exit(0)
	return
}

func (s *Keygen) Clean() {

}
