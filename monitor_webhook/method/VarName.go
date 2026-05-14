package method

type RequestBody struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type Result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type RobotInfo struct {
	routes map[string][]Notifier
}

func InitRobot(cfg *Config) *RobotInfo {
	info := &RobotInfo{
		routes: make(map[string][]Notifier),
	}
	for name, rc := range cfg.Routes {
		useDing := rc.UseDing == nil || *rc.UseDing
		useEmail := rc.UseEmail == nil || *rc.UseEmail

		var notifiers []Notifier
		if useDing && rc.DingTalk != "" {
			notifiers = append(notifiers, NewDingTalk(rc.DingTalk))
		}
		if useEmail && len(rc.Email) > 0 && cfg.SMTP.Host != "" {
			notifiers = append(notifiers, NewEmail(&cfg.SMTP, rc.Email))
		}
		info.routes[name] = notifiers
	}
	return info
}
