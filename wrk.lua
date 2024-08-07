math.randomseed(os.time()) -- Initialize random seed

request = function()
	local to_account_id = math.random(1, 1000)
	local from_account_id = math.random(1, 1000)
	local amount = math.random(1, 1000000)

	-- Ensure from_account_id is different from to_account_id
	while from_account_id == to_account_id do
		from_account_id = math.random(1, 1000)
	end

	local body = string.format(
		'{"to_account_id": %d, "from_account_id": %d, "amount": %d}',
		to_account_id,
		from_account_id,
		amount
	)

	local headers = {}
	headers["Content-Type"] = "application/json"

	return wrk.format("POST", "/transfers", headers, body)
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
