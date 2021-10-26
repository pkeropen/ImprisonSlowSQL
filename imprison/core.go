package imprison

import "ImprisonSlowSQL/pkg"

func verifyMySQL() {

	err := pkg.ExecCommand("/usr/sbin/getcap $(which mysql)")
	if err != nil {

	}
}
