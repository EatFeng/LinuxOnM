package routers

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&DashboardRouter{},
		&LogRouter{},
		&AuthRouter{},
		&TerminalRouter{},
		&FileRouter{},
		&HostRouter{},
		&GroupRouter{},
		&SettingRouter{},
		&CronjobRouter{},
		&ProcessRouter{},
	}
}
