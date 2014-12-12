helpdesk
.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: '<h1>Welcome to the helpdesk</h1>'
	})
	.when('/tickets/mine', {
		templateUrl: '/templates/ticketlist.html',
		controller: 'myTicketListCtrl'
	})
	.when('/tickets/department', {
		templateUrl: '/templates/ticketList.html',
		controller: 'depTicketListCtrl'
	})
	.when('/tickets/assigned', {
		templateUrl: '/templates/ticketList.html',
		controller: 'assignedTicketListCtrl'
	})
	.when('/ticket/:id', {
		templateUrl: 'ticket.html'
	})
});