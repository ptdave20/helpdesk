helpdesk
.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: '<h1>Administration</h1>'
	})
	.when('/domain', {
		template: '<h1>Domain Settings</h1>'
	})
	.when('/users', {
		template: '<h1>Users Settings</h1>'
	})
});