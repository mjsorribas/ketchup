import Theme from 'lib/theme';
import Layout from 'components/layout';

export default class ThemePage {
  theme: Mithril.Property<Theme>;

  constructor() {
    this.theme = m.prop<Theme>();
    let themeName = m.route.param('name');
    if (themeName) {
      Theme.get(themeName).then((theme) => {
        this.theme(theme);
      });
    }
  }

  static controller = ThemePage;
  static view(ctrl: ThemePage) {
    return Layout(
      !ctrl.theme() ? ''
        :
        m('.theme', [
          m('h1',
            m.trust('Theme &rsaquo; '),
            ctrl.theme().name
          ),
          m('h2', 'Templates'),
          m('table',
            Object.keys(ctrl.theme().templates).sort().map((name) => {
              let t = ctrl.theme().templates[name];
              return m('tr', [
                m('td', t.name),
                m('td', t.engine)
              ]);
            })
          ),
          m('h2', 'Assets'),
          m('table',
            Object.keys(ctrl.theme().assets).sort().map((a) => {
              return m('tr', [
                m('td', a),
              ]);
            })
          )
        ])
    );
  }
}