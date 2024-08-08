math.randomseed(os.time()) -- Initialize random seed

request = function()
	local to_account_id = math.random(1, 1000)
	local from_account_id = math.random(1, 1000)
	local amount = math.random(1, 1000000)

	-- Ensure from_account_id is different from to_account_id
	while from_account_id == to_account_id do
		from_account_id = math.random(1, 1000)
	end

	local json = string.format(
		[[
{
  "Description": "%s",
  "Entrys": [
    {
      "AccountID": %d,
      "Amount": %d,
      "Type": "%s"
    },
    {
      "AccountID": %d,
      "Amount": %d,
      "Type": "%s"
    }
  ]
}]],
		"Test Transaction", -- Description
		from_account_id,
		-amount,
		"DEBIT", -- First Entry
		to_account_id,
		amount,
		"CREDIT" -- Second Entry
	)

	local headers = {}
	headers["Content-Type"] = "application/json"

	return wrk.format("POST", "/transactions", headers, json)
end

responses = {}

response = function(status, headers, body)
	responses[status] = (responses[status] or 0) + 1
end

done = function()
	print("\nStatus code distribution:")
	for status, count in pairs(responses) do
		print(string.format("Status %d: %d responses", status, count))
	end
end
