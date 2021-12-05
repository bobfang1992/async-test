
from os import read
import time
import httpx
import asyncio
import csv
import json

BASE_URL = "https://pokeapi.co/api/v2/pokemon/"


async def get_pokemons(number: int):

    coroutines = []
    timeout = httpx.Timeout(60, read=30)
    async with httpx.AsyncClient(timeout=timeout) as client:
        for i in range(1, number+1):
            r = client.get(BASE_URL + str(i))
            coroutines.append(r)
        result = await asyncio.gather(*coroutines)

    return result


def main(start_number: int, end: int, step: int):
    async_time = []

    # async
    for number in range(start_number, end+1, step):
        start = time.time()
        result = asyncio.run(get_pokemons(number))
        print("Time to get {} pokemons: {} seconds".format(
            number, time.time() - start))
        async_time.append((number, time.time() - start))
        jsons = []
        for r in result:
            try:
                jsons.append(r.json())
            except json.decoder.JSONDecodeError:
                print(r.text)
        print(f"Got {len(jsons)} pokemons")
        print(f"First 10 pokemons: {[p['name'] for p in jsons[:10]]}")

    with open("async_time.csv", "w") as f:
        csv_out = csv.writer(f)
        csv_out.writerow(["number", "time"])
        for row in async_time:
            csv_out.writerow(row)

    # sync
    sync_time = []
    for number in range(start_number, end+1, step):
        start = time.time()
        result = []
        for i in range(1, number+1):
            r = httpx.get(BASE_URL + str(i))
            result.append(r)
        print("Time to get {} pokemons: {} seconds".format(
            number, time.time() - start))
        sync_time.append((number, time.time() - start))
        jsons = []
        for r in result:
            try:
                jsons.append(r.json())
            except json.decoder.JSONDecodeError:
                print(r.text)
        print(f"Got {len(jsons)} pokemons")
        print(f"First 10 pokemons: {[p['name'] for p in jsons[:10]]}")

    with open("sync_time.csv", "w") as f:
        csv_out = csv.writer(f)
        csv_out.writerow(["number", "time"])
        for row in sync_time:
            csv_out.writerow(row)


if __name__ == "__main__":
    main(10, 800, 20)
