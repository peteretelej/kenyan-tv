Vue.component('channel-listing',{
	props: ['channels'],
	template: '<div class="section">'+
	'<div class="container">'+
	'<div class="columns is-multiline">'+
	'<div v-for="channel in channels" class="column">'+
	'{{channel.Name}}'+
	'</div><!--/.column-->'+
	'</div>'+
	'</div><!--/.container-->'+
	'</div><!--/.section-->'
})

