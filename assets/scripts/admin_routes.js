admin
.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: '<h1>Administration</h1>'
	})
	.when('/domain', {
		template: '<h1>Settings</h1>'
	})
	.when('/buildings', {
		templateUrl: '/templates/admin_buildings.html',
		controller: 'bldgCtrl'
	})
	.when('/departments', {
		templateUrl: '/templates/admin_departments.html',
		controller: 'depCtrl'
	})
	.when('/department/config/:id', {
		templateUrl: '/templates/admin_config_dep.html',
		controller: 'depConfigCtrl'
	})
	.when('/department/tickets/:id', {
		template: '<h1>Show Department Tickets</h1>'
	})
	.when('/department/reports/:id', {
		template: '<h1>Show Department Reports</h1>'
	})
	.when('/users', {
		template: '<h1>Users</h1>'
	})
});