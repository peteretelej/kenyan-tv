var cacheName = "theoktv-v5"
self.addEventListener("install", function(event){
	event.waitUntil(
		caches.open(cacheName).then(function(cache){
			cache.addAll([
				"/",
				"/static/manifest.json",
				"/static/icons/apple-touch-icon.png",
				"/static/icons/favicon-32x32.png",
				"/static/icons/favicon-16x16.png",
				"/static/icons/browserconfig.xml",
				"/static/css/font-awesome/css/font-awesome.min.css",
				"/static/css/font-awesome/fonts/fontawesome-webfont.woff2?v=4.7.0",
				"/static/css/font-awesome/fonts/fontawesome-webfont.woff?v=4.7.0",
				"/static/css/font-awesome/fonts/fontawesome-webfont.ttf?v=4.7.0",
				"/static/js/home.js",
				"/static/js/main.js",
				"/static/css/bulma.min.css",
				"/static/js/vendor/axios.min.js",
				"/static/js/vendor/vue.js"
		])
		})
	)
})

self.addEventListener("activate",function(event){
	event.waitUntil(
		caches.keys().then(function(keys){
			keys.forEach(function(key){
				if (!key.startsWith("theoktv")){
					return
				}
				if (key !== cacheName){
					caches.delete(key)
				}
			})
		})
	)
})

self.addEventListener("fetch",function(event){
	event.respondWith(
		caches.open(cacheName).then(function(cache){
			return caches.match(event.request).then(function(response){
				return response ||fetch(event.request)
			})
		})
	)
})
