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

Vue.component("current-channel",{
	props: ['channel'],
	template: '<div>'+
	'<div style="display:block">'+
	'<iframe :key="channel.id" width="700px" height="500px" v-bind:src="channel.livestream" frameborder="0" allowfullscreen></iframe>'+
	'</div>'+
	'</div>'
} )
