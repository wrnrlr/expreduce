Rectangle[] := Rectangle[{0,0}]
Rectangle::usage = "`Rectangle[expr]` .";
Attributes[Rectangle] = {Protected};

Circle[{x_, y_}]:= Circle[{x,y}, 1]
Circle::usage = "`Circle[{x,y} r]` ."
Attributes[Circle] = {Protected};
