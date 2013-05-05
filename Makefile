all: frontend backend

backend: backend/*.go
	go build -o tasker -p 2 ./backend/

frontend: css js

css: static/main.css

static/main.css: frontend/less/index.less frontend/less/main.less \
				 frontend/less/color.less frontend/less/controls.less \
				 frontend/less/overlay.less
	lessc -O2 -x less/index.less static/main.css

#static/mobile.js
js: static/main.js static/desktop.js

static/main.js: frontend/js/main.js frontend/js/api.js
	@closure --js $^ > static/main.js

static/desktop.js: frontend/js/desktop.js
	@closure --js $^ > static/desktop.js

static/mobile.js:
	@echo "Mobile unimplemented."
