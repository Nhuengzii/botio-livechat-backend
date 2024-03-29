.PHONY: deploy init destroy

init:
	terraform init -backend-config="terraform_backend.conf"
	@echo "Initialized!"

deploy:
	terraform apply -auto-approve -var-file="terraform.tfvars"
	terraform apply -auto-approve -var-file="terraform.tfvars"
	@echo "Deployed!"
	
apply:
	terraform apply -auto-approve -var-file="terraform.tfvars"
	@echo "Applied!"

destroy:
	terraform destroy -auto-approve
	@echo "Destroyed!"
