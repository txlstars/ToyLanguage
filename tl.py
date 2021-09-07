#!/usr/bin/env python3

import argparse
import subprocess

args = None


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--gl", help="generate lexer source code by g4", action="store_true"
    )
    parser.add_argument(
        "--gp", help="generate parser source code by g4", action="store_true"
    )
    parser.add_argument(
        "--glp", help="generate lexer and parser source code by g4", action="store_true"
    )
    global args
    args = parser.parse_args()


antlr4 = "java -jar ./tools/antlr-4.9.2-complete.jar"


def generate_lexer():
    subprocess.run(
        "{} -Xexact-output-dir -o ./lexer -package lexer -Dlanguage=Go ./lexer/GoLexer.g4".format(
            antlr4
        ),
        shell=True,
        check=True,
        capture_output=True,
    )


def generate_parser():
    subprocess.run(
        "{} -Xexact-output-dir -o ./parser -package parser -Dlanguage=Go -lib ./lexer ./parser/GoParser.g4".format(
            antlr4
        ),
        shell=True,
        check=True,
        capture_output=True,
    )


def main():
    parse_args()

    if args.gl or args.glp:
        generate_lexer()

    if args.gp or args.glp:
        generate_parser()


if __name__ == "__main__":
    main()
