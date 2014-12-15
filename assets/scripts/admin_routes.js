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
		template: '<h1>Buildings</h1>'
	})
	.when('/departments', {
		templateUrl: '/templates/admin_departments.html',
		controller: 'depCtrl'
	})
	.when('/users', {
		template: '<h1>Users</h1>'
	})
});