# Steps
1. Change the env variables according to your db configuration in .env file

2. Make sure the table "settlement_logs" is added in the database 

3. Update the data of merchant Id and payment id in data.go file according to the following example
	```
	Example := 
	data = []IDs{
		{"mer_MsCtIPhqRc8045", "pay_sS3mEMr8ot2551"},
		{"mer_MsCtIPhqRc8045", "pay_sS3mEMr8ot2590"},
		                    .
		                    .
		                    .
		{"mer_MsCtIPhqRc8030", "pay_sS3mEMr8ot2578"},					
	}	
	```
4. Run `go run .` in terminal

5. Data is updated in "settlement_details" table and inserted details in "settlement_logs" table


# Conditions
1. If successfully updated
	```
   msg := updated for merchant id: <merchant_id> and payement id: <payment_id>
   example := updated for merchant id:  mer_MsCtIPhqRc8045  and payement id:  pay_sS3mEMr8ot2557
	```

2. If no row is present in the table
   ```
   msg := no rows to updated for merchant id: <merchant_id> and payment id: <payment_id>
   example := no rows to updated for merchant id: mer_MsCtIPhqRc8045 and payment id: pay_sS3mEMr8ot2590
	```
