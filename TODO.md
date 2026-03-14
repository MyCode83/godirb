# TODO – godirb

> CTRL + SHIFT + V

> Estado actual: DIR y FUZZ funcionales, base estable.
> Prioridad: consolidar fuzzing y UX antes de añadir features grandes.

---

##  Core / Motor

- [ ] Mejorar motor de fuzzing
  - [ ] Revisar baseline (status + length) -
  - [ ] Ajustar tolerancia por defecto -
  - [ ] Evitar falsos positivos obvios
  - [ ] Output coherente y consistente (`[~]`, `[+]`, etc.)
  - [ ] Asegurar que FUZZ nunca toca wildcard

- [ ] Port fuzzing
  - [ ] Definir semántica correcta (error ≠ inexistente)
  - [ ] Baseline específico para puertos
  - [ ] Sin wildcard
  - [ ] Output propio y claro
  - [ ] Integrar sin romper dir/fuzz actuales

- [ ] Query fuzzing
  - [ ] Placeholder en query (`?id=FUZZ`)
  - [ ] Reusar baseline existente
  - [ ] Confirmar compatibilidad con filtros actuales

---

##  UX / Errores

- [ ] Mejorar mensajes de error
  - [ ] URL inválida
  - [ ] Placeholder no encontrado
  - [ ] Flags incompatibles
  - [ ] Salidas limpias (sin panic)

- [ ] Revisar códigos de salida (exit codes)

---

##  CLI / Help

- [ ] Actualizar `--help`
  - [ ] Explicar modo DIR
  - [ ] Explicar modo FUZZ
  - [ ] Ejemplos claros y simples
  - [ ] Defaults visibles

---

##  Leaf-features (solo si no rompen core)

- [ ] Pequeñas flags opcionales
- [ ] Mejoras cosméticas
- [ ] Refinar output sin tocar lógica

---

##  Distribución

- [ ] README.md
  - [ ] Qué es godirb
  - [ ] Cuándo usar DIR vs FUZZ
  - [ ] Ejemplos básicos
  - [ ] Advertencias (port fuzzing, fuzz avanzado)

- [ ] install.sh
- [ ] install.ps1
- [ ] Preparar releases
  - [ ] Binarios
  - [ ] Changelog
  - [ ] Versionado limpio

---

##  Limpieza final

- [ ] Revisar logs de debug
- [ ] Eliminar código muerto
- [ ] Revisar nombres y consistencia

