GOLANG=go

LOG_LEVEL=$LOG_LEVEL $GOLANG run ./main.go \
  runOperation \
  -k ../aggspotter_keys/keeper1.json \
  -k ../aggspotter_keys/keeper2.json \
  -k ../aggspotter_keys/keeper3.json \
  -o 3XDRdtEnSqMb3zMyxobq3xt3Uq8oPcvGfAk4WkZXGhQ9
