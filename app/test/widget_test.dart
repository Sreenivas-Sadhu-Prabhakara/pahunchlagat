import 'package:flutter_test/flutter_test.dart';

import 'package:pahunchlagat_app/main.dart';

void main() {
  testWidgets('shows landed cost from defaults', (tester) async {
    await tester.pumpWidget(const PahunchLagatApp());
    await tester.tap(find.text('Landed cost'));
    await tester.pump();
    expect(find.text('Landed cost per unit'), findsOneWidget);
  });

  test('landed mirrors the Go oracle: composition adds GST', () {
    final r = landed(goods: 10000, gstPct: 18, freight: 500, coolie: 200,
        transport: 300, wastagePct: 5, units: 100, regular: false);
    expect(r.gstInCost, closeTo(1800, 1e-6));
    expect(r.landedTotal, closeTo(12800, 1e-6));
  });
}
