curl -k -d "@psmvpc.json" -X POST https://localhost:8089/psm/configs/network/v1/tenant/123/virtualrouters| python -m json.tool
