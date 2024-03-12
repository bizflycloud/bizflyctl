package constants

const (
	HaNoiRegion     = "HaNoi"
	HoChiMinhRegion = "HoChiMinh"
	VcHaNoiRegion   = "VC-HaNoi"
)

var RegionMapping = map[string]string{
	"hn":        HaNoiRegion,
	"hanoi":     HaNoiRegion,
	"hcm":       HoChiMinhRegion,
	"hochiminh": HoChiMinhRegion,
	"vc-hanoi":  VcHaNoiRegion,
}
