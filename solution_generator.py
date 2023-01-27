from __future__ import annotations
import itertools

import random
import yaml
import click


def connections_capacity(
    higher_order_entities: list[BaseEntity], lower_order_entities: list[BaseEntity]
) -> list[int]:
    return [
        (
            min_constraint := random.randint(
                1, int((max_constraint := min(hoe.max_capacity, loe.max_capacity))/10)
            ),
            random.randint(min_constraint, max_constraint),
        )
        for hoe in higher_order_entities
        for loe in lower_order_entities
    ]


def generate_mapping(
    suppliers: list[Supplier],
    factories: list[Factory],
    warehouses: list[Warehouse],
    shops: list[Shop],
) -> dict:
    mapping = {
        "D": len(suppliers),
        "F": len(factories),
        "M": len(warehouses),
        "S": len(shops),
        "sd": [e.max_capacity for e in suppliers],
        "sf": [e.max_capacity for e in factories],
        "sm": [e.max_capacity for e in warehouses],
        "ss": [e.max_capacity for e in shops],
        "cd": list(itertools.chain.from_iterable(
            [e.connection_price_to_lower for e in suppliers]
        )),
        "cf": list(itertools.chain.from_iterable(
            [e.connection_price_to_lower for e in factories]
        )),
        "cm": list(itertools.chain.from_iterable(
            [e.connection_price_to_lower for e in warehouses]
        )),
        "ud": [e.setup_cost for e in suppliers],
        "uf": [e.setup_cost for e in factories],
        "um": [e.setup_cost for e in warehouses],
        "p": [e.shop_income for e in shops],
        "xdminmax": list(itertools.chain.from_iterable(
            connections_capacity(suppliers, factories))
        ),
        "xfminmax": list(itertools.chain.from_iterable(
            connections_capacity(factories, warehouses))
        ),
        "xmminmax": list(itertools.chain.from_iterable(
            connections_capacity(warehouses, shops))
        ),
    }
    return mapping


class BaseEntity:
    def __init__(self, number_of_lower: int) -> None:
        self.max_capacity = random.randint(15, 100)
        self.setup_cost = random.randint(60, 150)
        self.connection_price_to_lower = [
            random.randint(1, 20) for _ in range(number_of_lower)
        ]


class Supplier(BaseEntity):
    ...


class Factory(BaseEntity):
    ...


class Warehouse(BaseEntity):
    ...


class Shop(BaseEntity):
    def __init__(self, number_of_lower: int) -> None:
        self.shop_income = random.randint(500, 1000)
        super().__init__(number_of_lower)


@click.command()
@click.argument("suppliers_count", type=click.INT)
@click.argument("factories_count", type=click.INT)
@click.argument("warehouses_count", type=click.INT)
@click.argument("shops_count", type=click.INT)
@click.argument("filename", type=click.STRING)
def main(suppliers_count, factories_count, warehouses_count, shops_count, filename):
    suppliers = [Supplier(factories_count) for _ in range(suppliers_count)]
    factories = [Factory(warehouses_count) for _ in range(factories_count)]
    warehouses = [Warehouse(shops_count) for _ in range(warehouses_count)]
    shops = [Shop(0) for _ in range(shops_count)]

    mapping = generate_mapping(suppliers, factories, warehouses, shops)

    print(mapping)
    print(filename)
    with open(f"solution_{filename}_{suppliers_count}_{factories_count}_{warehouses_count}_{shops_count}.yaml", "w") as file:
        file.write(yaml.safe_dump(mapping))

if __name__ == "__main__":
    main()
