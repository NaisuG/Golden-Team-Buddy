import json
import requests

URL = "https://tft.dakgg.io/api/v1/data/items"
PARAMS = {"hl": "en", "season": "set17"}


def main():
    response = requests.get(URL, params=PARAMS, timeout=10)
    response.raise_for_status()
    data = response.json()

    with open("raw_items.json", "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

    print("Tipo del JSON completo:", type(data))
    if isinstance(data, dict):
        print("Claves de primer nivel:", list(data.keys()))
    print("Total de items (crudo):", len(data.get("items", [])))


if __name__ == "__main__":
    main()