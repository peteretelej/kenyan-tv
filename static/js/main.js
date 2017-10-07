var vm = new Vue({
	el: "#app",
	data:{
		currentView:'home',
		currentChannel:'',
		offline:false,
		channels: [],
		intlchannels: []
	},
	created: function(){
		var self = this
		axios.get("/api?q=channels").then(function(response){
			if (response.data.code !== 200) return
			self.channels = response.data.data
			self.currentChannel = self.channels[0]
		}).catch(function(err){
			self.offline =true
		})
		axios.get("/api?q=intlchannels").then(function(response){
			if (response.data.code !== 200) return
			self.intlchannels = response.data.data
		}).catch(function(err){
			self.offline =true
		})
	},
	methods:{
		changeChannel: function(ch){
			this.currentChannel = ch
		}
	}

})
