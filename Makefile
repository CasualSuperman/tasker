.SUFFIXES: 

all: frontend backend

clean:
	@echo -e "CLEAN\ttasker"
	@rm -f tasker
	@echo -e "CLEAN\tmain.js"
	@rm -r static/main.js
	@echo -e "CLEAN\tdesktop.js"
	@rm -r static/desktop.js
	@echo -e "CLEAN\tresponsive.js"
	@rm -r static/responsive.js

backend: tasker

tasker: backend/*.go
	@echo -e "6g\ttasker"
	@go build -o tasker -p 2 ./backend/

frontend: css js

css: static/main.css

static/main.css: frontend/less/index.less frontend/less/main.less \
				 frontend/less/color.less frontend/less/controls.less \
				 frontend/less/overlay.less
	@echo -e "LESS\tindex.less"
	@lessc -O2 -x frontend/less/index.less static/main.css

#static/mobile.js
js: static/main.js static/desktop.js static/responsive.js

static/main.js: frontend/js/main.js frontend/js/api.js
	@echo -e "JS\tmain.js"
	@closure --js $^ > static/main.js

static/desktop.js: frontend/js/desktop.js
	@echo -e "JS\tdesktop.js"
	@closure --js $^ > static/desktop.js

static/responsive.js: frontend/js/responsive.js
	@echo -e "JS\tresponsive.js"
	@closure --js $^ > static/responsive.js

static/mobile.js:
	@echo "Mobile unimplemented."
