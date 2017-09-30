var vm = new Vue({
	el: "#app",
	data:{
		channels: [],
	},
	created: function(){
		var self = this
		axios.get("/api?q=channels").then(function(response){
			if (response.data.code !== 200) return
			self.channels = response.data.data
		})
	}

})
