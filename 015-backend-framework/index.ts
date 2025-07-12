class TrieNode {
    children: Map<string, TrieNode>;
    isEOW: boolean; // end of word

    constructor() {
        this.children  = new Map();
        this.isEOW = false;
    }

    add(letter: string, isLastChar: boolean): TrieNode {
        const newNode = new TrieNode();
        if (isLastChar) newNode.isEOW = true;

        this.children.set(letter, newNode);
        return newNode;
    }
}

class Trie {
    root: TrieNode;

    constructor() {
        this.root = new TrieNode();
    }

    insertWithRecursion(word: string, node: TrieNode = this.root) {
        if (word.length == 0) return;

        const foundNode = node.children.get(word[0]);
        if (foundNode) {
            this.insertWithRecursion(word.slice(1), foundNode);
            return;
        }

        const newNode = node.add(word[0], word.length == 1);
        this.insertWithRecursion(word.slice(1), newNode);
        return;
    }

    insert(word: string) {
        if (word.length == 0) return;

        let cur = this.root;
        for(let i=0; i < word.length; i++) {
            const curWord = word[i];
            if (!cur.children.has(curWord)) {
                cur.children.set(curWord, new TrieNode());
            }
            cur = cur.children.get(curWord)!;
        }
        cur.isEOW = true;
    }
}

const trie = new Trie()
trie.insert("hello")
trie.insert("he")
trie.insert("ha")
console.log(trie);
