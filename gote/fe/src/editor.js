require.config({ paths: { 'vs': 'https://cdn.jsdelivr.net/npm/monaco-editor@0.20.0/min/vs' }});

let monacoResolve;
const monacoReady = new Promise((resolve, reject) => {
  monacoResolve = resolve
})

require(['vs/editor/editor.main'], function() {
  monacoResolve(monaco)
});

// TODO: remove this
let languages = []

export const getmonaco = () => monacoReady

export async function mountEditor(el) {
  const monaco = await getmonaco()
  languages = monaco.languages.getLanguages()
  // TODO: destory editor
  console.log("create editor")
  const editor = monaco.editor.create(el, {
    scrollBeyondLastLine: false,
    automaticLayout: true
  });
  return editor
}

export function setText(editor, content, name='') {
  if (content === '') {
    editor.setValue('');
    return;
  }
  const lang = filename2language(name);
  editor.setModel(monaco.editor.createModel(content, lang));
  editor.layout();
};

export function getText(editor) {
  return editor.getValue()
}

function filename2language(name) {
  name = name.substr(name.lastIndexOf('/') + 1);
  const index = name.lastIndexOf('.');
  if (index === -1) {
    return '';
  }
  const ext = name.substr(index).toLowerCase();
  for (const lang of languages) {
    if (lang.extensions && lang.extensions.indexOf(ext) !== -1) {
      return lang.id;
    }
  }
}