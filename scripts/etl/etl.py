import argparse
import json
from pathlib import Path

import requests

SEASON = "set17"
LANGUAGE = "en"
API_URL = "https://tft.dakgg.io/api/v1/data/{dataset}"

SCRIPT_DIR = Path(__file__).resolve().parent
OUTPUT_DIR = SCRIPT_DIR.parent.parent / "internal" / "catalog" / "data"

CHAMPION_FIELDS = ["key", "ingameKey", "name", "imageUrl", "traits", "role"]
TRAIT_FIELDS = ["key", "ingameKey", "name", "imageUrl", "type", "styles", "stats"]
ITEM_FIELDS = ["key", "ingameKey", "name", "desc", "fromDesc", "imageUrl", "compositions", "affectedTraitKey"]


def raw_path(dataset):
    return SCRIPT_DIR / f"raw_{dataset}.json"


def read_json(path):
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)


def write_json(path, data):
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)


def download(dataset):
    response = requests.get(
        API_URL.format(dataset=dataset),
        params={"hl": LANGUAGE, "season": SEASON},
        timeout=10,
    )
    response.raise_for_status()
    write_json(raw_path(dataset), response.json())


def pick(raw, fields):
    return {field: raw.get(field) for field in fields}


def clean_champion(raw):
    cleaned = pick(raw, CHAMPION_FIELDS)
    cleaned["cost"] = raw["cost"][0]
    return cleaned


def clean_trait(raw):
    return pick(raw, TRAIT_FIELDS)


def item_category(raw):
    if raw.get("isFromItem"):
        return "component"
    if raw.get("isRadiant"):
        return "radiant"
    if raw.get("isEmblem"):
        return "emblem"
    if raw.get("isArtifact"):
        return "artifact"
    if raw.get("isNormal"):
        return "completed"
    if raw.get("isTrait") or raw.get("isPsionic") or raw.get("isAnimaSquad") or raw.get("isCommandMod"):
        return "special"
    return "other"


def clean_item(raw):
    cleaned = pick(raw, ITEM_FIELDS)
    cleaned["category"] = item_category(raw)
    return cleaned


CLEANERS = {
    "champions": clean_champion,
    "traits": clean_trait,
    "items": clean_item,
}


def clean(dataset):
    data = read_json(raw_path(dataset))
    if data["season"] != SEASON:
        raise ValueError(f"{raw_path(dataset).name} es de {data['season']}, se esperaba {SEASON}")

    visible = [entry for entry in data[dataset] if not entry.get("isHidden")]
    output = {"season": SEASON, dataset: [CLEANERS[dataset](entry) for entry in visible]}
    write_json(OUTPUT_DIR / f"{dataset}.json", output)
    print(f"{dataset}: {len(visible)} de {len(data[dataset])}")


def main():
    parser = argparse.ArgumentParser(description=f"Descarga y limpia los datos de TFT ({SEASON}).")
    parser.add_argument("--offline", action="store_true", help="limpia los raw_*.json existentes sin descargar")
    args = parser.parse_args()

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    for dataset in CLEANERS:
        if not args.offline:
            download(dataset)
        clean(dataset)


if __name__ == "__main__":
    main()
