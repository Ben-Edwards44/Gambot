#!/bin/bash

wget https://database.lichess.org/standard/lichess_db_standard_rated_2013-01.pgn.zst  #download chess game from the database here: https://database.lichess.org/

unzstd -d lichess_db_standard_rated_2013-01.pgn.zst  #unpack the file

rm lichess_db_standard_rated_2013-01.pgn.zst  #delete the compressed file (no longer needed)

#run the python scripts
python extract_fen.py
python get_equal.py