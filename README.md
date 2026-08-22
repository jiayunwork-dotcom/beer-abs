# beer-abs — 溶液吸光度核算

beer-abs 是一个命令行溶液吸光度核算工具。你输入摩尔吸光系数 ε、浓度 c、光程 L（以及可选的杂散光分数 s 与有限带宽参数），它按 Beer–Lambert 定律算出吸光度 A 与透射比 T；多组分按无相互作用可加求和，并给出杂散光修正后的观测值 A_obs、T_obs。能力边界：Beer–Lambert 理想模型加上杂散光修正，组分之间假定无相互作用；不做配液指导，也不做色谱塔板高度等无关核算。

- 输入：JSON 算例文件（`components` 数组 + `path_length`，可选 `stray_fraction` 与 `band`），示例见 `example/cu-508nm.json`、`example/mix-cu-ni.json`、`example/cu-508nm-band.json`
- 输出：各组分的 A_i、理想总吸光度 A 与透射比 T、杂散光修正后的 A_obs 与 T_obs、偏离量 deviation
- 非法输入（ε≤0、c<0、L≤0、s 不在 [0,1)、未知 JSON 字段）一律报 stderr 并非零退出

## 用法

在仓库根目录运行：

```text
go run . absorbance example/cu-508nm.json
```

输出示例（A=εcL 闭式值 1.0、T=10⁻¹）：

```text
sample: 1 component(s), path L = 1 cm, stray s = 0
  [0] Cu: eps = 125, c = 0.008, A_i = 1
A = 1   T = 0.1
A_obs = 1   T_obs = 0.1   s = 0
deviation = +0  (ideal instrument)
```

算例文件字段：

| 字段 | 含义 | 必填 |
|------|------|------|
| `components[].label` | 组分名 | 否 |
| `components[].extinction` | 摩尔吸光系数 ε | 是（无 `band` 时） |
| `components[].extinction_low` / `extinction_high` | 带宽两端波长处的 ε | 是（有 `band` 时） |
| `components[].concentration` | 浓度 c | 是 |
| `path_length` | 光程 L (cm) | 是 |
| `stray_fraction` | 杂散光分数 s，默认 0 | 否 |
| `band.center` / `band.half_width` | 矩形通带中心与半宽 (nm) | 否 |

## 关键约定

- **单组分**：A = ε·c·L，要求 ε>0、c≥0、L>0。
- **透射比**：T = 10^(−A)，对数一律以 10 为底（等价地 T = e^(−ln10·A)），不允许与 log10 混用别的底。
- **多组分可加**：无相互作用时 A_tot = Σ ε_i·c_i·L，所有组分共用同一个 L。
- **杂散光修正**：T_obs = (T + s)/(1 + s)，A_obs = −log10(T_obs)；s ∈ [0,1)。s=0 时退回理想值；s>0 时高吸光度处 A_obs 低于 A_ideal（负偏离），低吸光度处接近理想。
- **透明样品**：c=0 时 A=0、T=1，与 s 无关。
- **倍增规则**：只把 c（或 L）加倍，A 加倍、T 变为原来的平方。
- **有限带宽**：带 `band` 时用矩形通带两端波长的 ε 平均值作为有效 ε，并报告与单波长读数的偏离。

## 构建与测试

```text
go build ./...       # 编译（纯标准库，无第三方依赖）
go test ./...        # 全部单元测试（law / band / mixture）
```

## 许可

MIT，见 [LICENSE](./LICENSE)。
