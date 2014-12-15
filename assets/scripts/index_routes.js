indexApp
.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: '/templates/index.html',
		controller: 'indexRoute'
	})
	.when('/error/invalid_domain', {
		templateUrl: '/templates/index.html',
		controller: 'invalidDomain'
	})
});