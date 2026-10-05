# Golden Team Buddy

App de escritorio para Teamfight Tactics (Set 17). A partir de los campeones que tienes en tablero y banca, sugiere hasta tres composiciones hacia las que avanzar, cada una con una posible sub-variante.

> **Estado: en pausa.** El proyecto se quedó en el Set 17. Con la llegada de la nueva temporada cambiaron los campeones y los traits, y no seguí con el proyecto porque esa temporada no me gustó. Planeo retomarlo a futuro.

Este repositorio es el primer prototipo (v1) de una versión mejorada que tengo planeada. Lo publico como parte de mi portafolio.

## Stack

- Go + [Wails v2](https://wails.io)
- JavaScript (Vite) para la interfaz
- Python para el ETL de datos (fuente: dakgg)

## Cómo funciona

El motor (`internal/engine`) arma las composiciones con un beam search. Cada línea se compara en cascada:

1. Fuerza de traits: breakpoints cruzados.
2. Campeones del tablero que aportan a un trait activo (la banca no tiene preferencia, solo cuenta como gratis).
3. Sinergia promedio entre pares de campeones.

Solo se consideran campeones que pueden aparecer en tienda al nivel actual.

## Estructura

```
app.go, main.go        Entrada de Wails y bindings con el frontend
internal/catalog       Modelo de datos y carga del catálogo embebido
internal/engine        Probabilidades, puntaje y búsqueda
internal/board         Tablero y banca
frontend/              Interfaz
scripts/etl/           Descarga y limpieza de datos (etl.py)
cmd/checkdata/         Prueba rápida del motor por consola
```

## Desarrollo

```bash
wails dev
```

```bash
go test ./...
```

```bash
wails build
```

Actualizar los datos del set (requiere `pip install -r scripts/etl/requirements.txt`):

```bash
python scripts/etl/etl.py
```

## Próximas versiones

### v1.2: cómo se obtienen y procesan los datos

El motor funciona, pero todavía no hace algunas distinciones importantes:

- Lo que ya tienes en el tablero pesa muy poco frente a la fuerza de traits, así que muchas veces sugiere cambiar casi todo.
- El puntaje de un trait depende de cuántos breakpoints tiene y no de su color (bronce, plata, oro, cromático).
- El nivel se calcula a partir de la cantidad de campeones en el tablero, en vez de ser un dato propio.
- Traits como Stargazer, que cambian en cada partida, usan breakpoints genéricos.

Ya estoy evaluando distintas soluciones para cada punto.

### v2: pruebas con datos reales

Probar el motor con resultados de partidas reales y usar esos datos para calibrar y mejorar el módulo completo.

## Datos

Los datos de campeones y traits se obtienen de la API pública de [dakgg](https://tft.dakgg.io), que no es una fuente oficial de Riot Games. Pueden contener errores o estar desactualizados, y pertenecen a sus respectivos dueños. El repositorio solo incluye el extracto reducido que la app necesita para funcionar. Es posible que la API ya no entregue los datos del Set 17, por lo que el script de actualización podría no funcionar.

## Aviso legal

Este es un proyecto personal, sin fines comerciales, hecho solo para mostrar mi trabajo. No está afiliado, patrocinado ni respaldado por Riot Games, y no busca infringir ninguno de sus derechos. Teamfight Tactics y todos los nombres, campeones y recursos asociados son propiedad de Riot Games, Inc.

Golden Team Buddy isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games, and all associated properties are trademarks or registered trademarks of Riot Games, Inc.
