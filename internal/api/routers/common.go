package routers

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&DashboardRouter{},
	}
}
