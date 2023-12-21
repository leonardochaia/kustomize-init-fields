
build: 
	docker build -t lchaia/kustomize-init-fields .

test: build
	kubectl kustomize --enable-alpha-plugins ./test