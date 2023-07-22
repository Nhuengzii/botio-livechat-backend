.PHONY: deploy init destroy

init:
	terraform init
	@echo "Initialized!"

deploy:
	terraform apply -auto-approve -var-file="terraform.tfvars"
	terraform apply -auto-approve -var-file="terraform.tfvars"
	@echo "Deployed!"
	
destroy:
	terraform destroy -auto-approve
	@echo "Destroyed!"
