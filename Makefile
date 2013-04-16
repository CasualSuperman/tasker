css: static/main.css

static/main.css: less/main.less
	lessc -O2 -x less/main.less static/main.css

js: base.js desktop.js mobile.js

base.js: js/main.js js/api.js
	@echo "main.js"
	@closure --js $^ > static/main.js

desktop.js: js/desktop.js
	@echo "desktop.js"
	@closure --js $^ > static/desktop.js

mobile.js:
	@echo "Mobile unimplemented."
