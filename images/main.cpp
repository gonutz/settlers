#include "binpack.h"
#include <iostream>
using namespace std;

ostream& operator << (ostream& os, RectangleBinPack::Node* node) {
	os << "(" << node->x << ", " << node->y << ", " << node->width << ", " << node->height << ")" << endl;
	return os;
}

int main() {
	RectangleBinPack packer;
	packer.Init(1024, 1024);
	auto node = packer.Insert(80, 60);
	cout << node << endl;
	auto node2 = packer.Insert(80, 60);
	cout << node2 << endl;
	getchar();
	return 0;
}