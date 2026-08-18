import 'package:flutter/material.dart';

void main() => runApp(const PahunchLagatApp());

/// Pahunchlagat — delivered cost. Nets freight, coolie, transport, wastage and
/// the decisive GST treatment into a true landed cost per unit.
class PahunchLagatApp extends StatelessWidget {
  const PahunchLagatApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Pahunchlagat',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF2E5E9A), useMaterial3: true),
        home: const HomePage(),
      );
}

class Result {
  final double gstInCost, landedTotal, usableUnits, perUnit;
  const Result(this.gstInCost, this.landedTotal, this.usableUnits, this.perUnit);
}

/// landed mirrors backend/cost.go exactly.
Result landed({
  required double goods, required double gstPct, required double freight,
  required double coolie, required double transport, required double wastagePct,
  required double units, required bool regular,
}) {
  final gstInCost = regular ? 0.0 : goods * gstPct / 100;
  final total = goods + gstInCost + freight + coolie + transport;
  final usable = units * (1 - wastagePct / 100);
  final perUnit = usable > 0 ? total / usable : 0.0;
  return Result(gstInCost, total, usable, perUnit);
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _goods = TextEditingController(text: '10000');
  final _gst = TextEditingController(text: '18');
  final _freight = TextEditingController(text: '500');
  final _coolie = TextEditingController(text: '200');
  final _transport = TextEditingController(text: '300');
  final _wastage = TextEditingController(text: '5');
  final _units = TextEditingController(text: '100');
  bool _regular = false;
  Result? _r;

  double _n(TextEditingController c) => double.tryParse(c.text.trim()) ?? 0;
  void _calc() => setState(() => _r = landed(
        goods: _n(_goods), gstPct: _n(_gst), freight: _n(_freight),
        coolie: _n(_coolie), transport: _n(_transport), wastagePct: _n(_wastage),
        units: _n(_units), regular: _regular,
      ));

  @override
  Widget build(BuildContext context) {
    String money(double v) => '₹${v.toStringAsFixed(2)}';
    return Scaffold(
      appBar: AppBar(
        title: const Text('Pahunchlagat · landed cost'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        const Text('What does this stock really cost, delivered?',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
        const SizedBox(height: 12),
        _f(_goods, 'Goods value pre-GST (₹)'),
        Row(children: [Expanded(child: _f(_gst, 'GST %')), const SizedBox(width: 12), Expanded(child: _f(_units, 'Units'))]),
        Row(children: [Expanded(child: _f(_freight, 'Freight ₹')), const SizedBox(width: 12), Expanded(child: _f(_coolie, 'Coolie ₹'))]),
        Row(children: [Expanded(child: _f(_transport, 'Local transport ₹')), const SizedBox(width: 12), Expanded(child: _f(_wastage, 'Wastage %'))]),
        SwitchListTile(
          title: const Text('GST-registered (recover input tax)'),
          subtitle: const Text('Off = composition / unregistered — GST becomes a cost'),
          value: _regular,
          onChanged: (v) => setState(() { _regular = v; _calc(); }),
          contentPadding: EdgeInsets.zero,
        ),
        FilledButton.icon(onPressed: _calc, icon: const Icon(Icons.local_shipping), label: const Text('Landed cost')),
        const SizedBox(height: 20),
        if (_r != null) _card(money),
      ]),
    );
  }

  Widget _f(TextEditingController c, String label) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: TextField(controller: c,
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
          decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
          onChanged: (_) => _calc()),
      );

  Widget _card(String Function(double) money) {
    final r = _r!;
    return Card(
      color: Theme.of(context).colorScheme.primaryContainer,
      child: Padding(padding: const EdgeInsets.all(16),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          const Text('Landed cost per unit'),
          Text(money(r.perUnit), style: const TextStyle(fontSize: 32, fontWeight: FontWeight.bold)),
          const Divider(height: 20),
          _row('GST in cost', money(r.gstInCost)),
          _row('Landed total', money(r.landedTotal)),
          _row('Usable units', r.usableUnits.toStringAsFixed(1)),
        ])),
    );
  }

  Widget _row(String k, String v) => Padding(padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [Text(k), Text(v, style: const TextStyle(fontWeight: FontWeight.w600))]));
}
