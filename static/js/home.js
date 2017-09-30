Vue.component('channel-listing',{
	props: ['channels'],
	template: '<div class="section">'+
	'<div class="container">'+
	'<div class="columns is-multiline">'+
	'<div v-for="channel in channels" class="column is-narrow">'+
	'<div class="box">'+
	'{{channel.name}}'+
	'<div style="display:block"><iframe width="200" height="200" v-bind:src="channel.livestream" frameborder="0" allowfullscreen></iframe></div>'+
	'</div><!--/.box-->'+
	'</div><!--/.column-->'+
	'</div>'+
	'</div><!--/.container-->'+
	'</div><!--/.section-->'
})

