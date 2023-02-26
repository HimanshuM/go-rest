package builder

var allRoutes []*Route

func AllRoutes() []*Route {
	return allRoutes
}

func buildRoute(route *Route) {

}

func BuildRoutes() {
	for _, route := range allRoutes {
		buildRoute(route)
	}
}
