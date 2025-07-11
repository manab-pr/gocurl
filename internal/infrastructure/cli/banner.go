package cli

import "github.com/fatih/color"

func ShowBanner() {
	banner := `
   __  __    _    _   _    _    ____  
  |  \/  |  / \  | \ | |  / \  | __ ) 
  | |\/| | / _ \ |  \| | / _ \ |  _ \ 
  | |  | |/ ___ \| |\  |/ ___ \| |_) |
  |_|  |_/_/   \_\_| \_/_/   \_\____/ 

`
	color.HiMagenta(banner)
}
