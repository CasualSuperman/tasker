css: static/main.css

static/main.css: less/index.less less/main.less less/color.less less/controls.less less/overlay.less
	lessc -O2 -x less/index.less static/main.css

js: base.js desktop.js mobile.js

base.js: js/main.js js/api.js js/ui.js
	@echo "main.js"
	@closure --js $^ > static/main.js

desktop.js: js/desktop.js
	@echo "desktop.js"
	@closure --js $^ > static/desktop.js

mobile.js:
	@echo "Mobile unimplemented."
