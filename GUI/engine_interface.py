import subprocess


ENGINE_PATH = "C:\\Users\\Ben Edwards\\Documents\\Programming\\Python\\Projects\\chess-engine\\chess-engine\\chess-engine.exe"


class Engine:
    def __init__(self):
        self.process = subprocess.Popen(ENGINE_PATH, stdin=subprocess.PIPE, stdout=subprocess.PIPE)

        self.check_uci()

    def __del__(self):
        self.process.stdout.close()
        self.process.stdin.close()
        self.process.kill()

    def set_pos(self, move_list):
        #set the position to the startpos + the moves played in move_list
        moves = " ".join(move_list)
        self.send_cmd(f"position startpos moves {moves}")

    def get_move(self, **kwargs):
        #assumes the position has been set
        self.send_args("go", kwargs)

        output = ""
        while len(output) < 8 or output[:8] != "bestmove":
            output = self.read_line()

        output = output.split(" ")
        move = output[1]  #output will look like: bestmove e2e4 ponder c7c5

        return move

    def check_uci(self):
        self.send_cmd("uci")

        output = ""
        while output != "uciok":
            output = self.read_line()

            if output != "" and output[:2] != "id" and output[:6] != "option" and output != "uciok":
                raise Exception(f"UCI handshake failed")

    def check_ready(self):
        self.send_cmd("isready")

        output = self.read_line()

        if output != "readyok":
            raise Exception("Ready check failed")

    def read_line(self):
        #read a line of stdout

        if self.process.poll() is not None:
            raise Exception("UCI engine process finished")  #process has exited for some reason
        
        output = self.process.stdout.readline()
        text = output.decode().strip()

        return text
    
    def send_args(self, cmd_name, kwarg_dict):
        #send args in the UCI required format (name1 value1 name2 value2...)
        args = ""
        for k, v in kwarg_dict.items():
            args += f" {k} {v}"

        self.send_cmd(f"{cmd_name}{args}")
        
    def send_cmd(self, cmd):
        print(f"Sending: {cmd}")
        b_cmd = bytes(f"{cmd}\n", encoding="utf-8")

        self.process.stdin.write(b_cmd)
        self.process.stdin.flush()