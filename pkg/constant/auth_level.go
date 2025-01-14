package constant

const (
	Level_User   string = "USER"
	Level_Vendor string = "VENDOR"
)

func Authlevel() []string {
	return []string{
		Level_User,
		Level_Vendor,
	}
}
