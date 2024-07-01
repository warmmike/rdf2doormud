#!/usr/bin/env python
import rdflib
import io
import pydotplus
from IPython.display import display, Image
from rdflib.tools.rdf2dot import rdf2dot

project = "1-1"

def visualize(g):
    stream = io.StringIO()
    rdf2dot(g, stream, opts = {display})
    dg = pydotplus.graph_from_dot_data(stream.getvalue())
    dg.rankdir = "LR"
    dg.write_png(project + '.png')
    dg.write_pdf(project + '.pdf')
    png = dg.create_png()

def main():
    g = rdflib.Graph()
    ret = None
    try:
        with open('./'+project+'.ttl') as f: ret = f.readlines()
    except: pass
    if ret is not None:
        ret = ''.join(ret)

        g.parse(data=ret, format='turtle')

        visualize(g)
    return 0

if __name__ == "__main__":
    main()