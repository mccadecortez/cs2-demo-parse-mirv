# Usage

```bash
go run bin/parse_demo ... --demo-file example.dem
```

```log
Wrote json output to: output.json
Wrote xml output to: ..._output.xml
```

---

```bash
go run bin/mirv/mirv.go ... --demo-json output.json
```

mirv hosts a websocket server that sends the demo-tick, button pressed, view angle, and velocity
