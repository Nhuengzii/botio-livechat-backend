.PHONY: deploy init destroy

init:
	terraform init
	@echo "Initialized!"

deploy:
	terraform apply -auto-approve
	@([ $$? -ne 0 ]) || (echo "Deployment failed! will deploy again" && terraform apply -auto-approve])	
	@echo "Deployed!"
	
destroy:
	terraform destroy -auto-approve
	@echo "Destroyed!"
