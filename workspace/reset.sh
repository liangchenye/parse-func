rm -fr openssl-1.0.2k
tar xvf openssl-1.0.2k-hobbled.tar
cp ec_curve.c openssl-1.0.2k/crypto/ec
cp ectest.c openssl-1.0.2k/crypto/ec
cp m2.sh openssl-1.0.2k

cd openssl-1.0.2k
sh m2.sh > ../parse.log

