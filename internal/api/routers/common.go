package routers

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&DashboardRouter{},
		&LogRouter{},
		&AuthRouter{},
	}
}
